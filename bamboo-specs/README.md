# GitHub Actions -> Bamboo, mapped against this repo

`bamboo.yaml` in this folder is a working translation of `.github/workflows/ci.yaml` and
`.github/workflows/deployment.yaml`. This file explains the concepts behind it.

## The mental model shift

GitHub Actions is **one runtime**: a workflow file describes triggers, a DAG of jobs, and
each job downloads its own toolchain at runtime.

Bamboo splits the same thing across **four objects**, and only some of them are code:

| Object | Contains | Config as code? |
| --- | --- | --- |
| Project | grouping + permissions + shared vars | No (create in UI first) |
| Linked repository | clone URL, credentials, branch detection | No (create in UI first) |
| Plan (stages -> jobs -> tasks) | build, verify, package | Yes |
| Deployment project (environments -> tasks) | releases, promotion, rollback | Yes |
| Agents + capabilities | the installed toolchain | No (per-agent config) |

That last row is the one that trips people up coming from Actions. There is no
`actions/setup-go`. Either the agent already has Go (a *capability* the job declares as a
*requirement*), or the job runs inside a container.

## Direct equivalents

| GitHub Actions | Bamboo | Notes |
| --- | --- | --- |
| workflow file | plan | Identified by `project-key` + `key`, not filename |
| `job` | job | Jobs in a stage run in **parallel** |
| `step` | task | Sequential, stop on first failure |
| `needs:` | **stage boundary** | No per-job DAG. See below |
| `if: always()` | `final-tasks:` | Also runs when the build is manually stopped |
| `if: <condition>` | task `conditions:` | Only variable comparisons, no expression language |
| `runs-on: ubuntu-latest` | `requirements:` + `docker:` | Requirement matches an agent capability |
| `runs-on: self-hosted` | custom capability, e.g. `agent_type: vps` | Same idea, explicit |
| `uses: some/action@v1` | nothing | Rewrite as a script task, or use a Bamboo plugin |
| `actions/checkout` | `checkout` task | Added implicitly as task 1 if you omit it |
| `actions/setup-go` | agent capability, or `docker: golang:1.x` | Pick one per job |
| `actions/cache` | agent working dir + `clean` task | Bamboo reuses the workspace by default; caching is "don't clean" |
| `outputs:` (job/workflow) | `inject-variables` with `scope: RESULT` | Write `key=value` to a file, then inject |
| `env:` | task `environment:` | Space-separated, quote values with spaces |
| `secrets.FOO` | plan/environment variable, encrypted | Sensitive keys must be `BAMSCRT@...` or import fails |
| `${{ github.sha }}` | `${bamboo.planRepository.revision}` | Or `$bamboo_planRepository_revision` in scripts |
| `${{ github.run_id }}` | `${bamboo.buildResultKey}` | e.g. `API-CI-42` |
| `upload/download-artifact` | `artifacts:` + `artifact-subscriptions:` | Shared artifacts auto-subscribe in later stages |
| `on: push` | `triggers: [remote]` | Bitbucket DC pushes to Bamboo; no webhook to manage |
| `on: pull_request` | `branches: { create: for-pull-request }` | Creates a plan branch per PR |
| `on: schedule` | `triggers: [cron: ...]` | Quartz cron (6-7 fields, `?` allowed) |
| `workflow_dispatch` | manual stage / manual deployment | See "redeploy" below |
| `workflow_call` (reusable) | `!include 'file.yaml'` or YAML anchors | Textual include, not a call |
| `environment: production` | deployment environment + `deploy` permission | The reviewer gate is a permission, not a setting |
| deployment concept | release + environment | Bamboo has real first-class releases |
| test summary in UI | `test-parser` task | Without it, Bamboo only sees pass/fail |
| Jira link | `branches: { link-to-jira: true }` | Branch name with issue key -> build status on the issue |

## The four things that will actually bite you

### 1. `needs:` does not exist

Your CI has `vet`, `lint`, `test`, `format` all with `needs: extract-tools-version`. In
Bamboo that becomes two stages:

```
Stage "Prepare"  ->  Extract Tools Version
Stage "Verify"   ->  Vet | Lint | Test | Format   (parallel)
```

Consequence: a stage is a full barrier. If you had `A -> C` and `B -> D` as two independent
chains, Actions runs them concurrently; Bamboo forces `A,B` then `C,D`. You cannot express
partial ordering. If it really matters, split into two plans and wire them with
`dependencies:`.

### 2. Job outputs are files, not return values

`extract-tools-version` returning `go-version` becomes: write `tools.properties`, then
`inject-variables` with `scope: RESULT`. Scope is the whole trick:

- `LOCAL` — gone when the job ends
- `RESULT` — persisted on the build result, readable by later stages **and by the linked
  deployment project**

Reference it as `${bamboo.tools.go_version}` in config fields, or `$bamboo_tools_go_version`
inside a script (dots become underscores).

Note the ordering constraint this creates: the `Lint` job's container image is
`golang:${bamboo.tools.go_version}-alpine`, so it *must* be in a stage after the job that
injects the variable. The extract job itself therefore can't run in that image.

### 3. Deploying is a separate object, and it's better than what you have

Your `deployment.yaml` reimplements release management by hand — `workflow_dispatch` with a
`sha` input, an image-existence check to make redeploys idempotent. Bamboo gives you:

- a **release**: an immutable named snapshot of one successful build result
- **deployment history** per environment, so "what is on prod" is a fact, not a guess
- **redeploy / rollback**: pick any past release, click deploy

So the `inputs.sha` escape hatch disappears. You don't retype a SHA, you select a release.

### 4. No service containers, no Docker CLI task in YAML

Your `vet` job spins up throwaway Postgres. Bamboo YAML Specs has:

- `docker:` on a job — runs **the job itself** in a container, not a sidecar
- no service containers
- no Docker CLI task (that exists in **Java** Specs only)

So the `docker run postgres ... && wait for healthy` script stays a script. Same for
`buildx build --push`. That part of the translation is nearly verbatim, which is convenient.

## Bamboo variables worth memorising

| Variable | Meaning |
| --- | --- |
| `${bamboo.buildNumber}` | incrementing build number |
| `${bamboo.buildResultKey}` | `PROJ-PLAN-42` |
| `${bamboo.planRepository.revision}` | commit SHA built |
| `${bamboo.planRepository.branch}` | branch name |
| `${bamboo.deploy.release}` | release name (deployments only) |
| `${bamboo.deploy.environment}` | environment name (deployments only) |
| `${bamboo.working.directory}` | job workspace |
| `${system.FOO}` | agent system/env variable |

In config fields use `${bamboo.x.y}`. In script bodies use `$bamboo_x_y`.

## Getting it running

1. Create the project in Bamboo (`API`). YAML Specs cannot create projects.
2. Create the linked repository named `api`, pointing at the Bitbucket repo.
3. On that linked repository, enable **Bamboo Specs** and grant the project RSS access.
4. Define agent capabilities the jobs require: `golang`, `hasDocker`, `agent_type=vps`.
5. Push to the default branch. Bamboo detects `bamboo-specs/bamboo.yaml` and creates the
   plan, deployment project and permissions. **The plan is now read-only in the UI** —
   edits happen only via this file.
6. Replace every `BAMSCRT@0@0@REPLACE_WITH_ENCRYPTED_VALUE` with a real encrypted value, or
   delete those keys and set them in the UI.

Debugging tip: when the import fails, Bamboo reports the failure on the repository's
**Bamboo Specs** tab, not on the plan. Errors look like
`Push to Test / tasks / [0] / inject-variables: Property is required.`

Reverse direction, which is the faster way to learn: build a plan by clicking in the UI,
then **Actions > View plan as YAML**. That is how you find the `any-task` plugin keys for
plugin tasks that Specs doesn't support natively.

## One note on Bitbucket Pipelines

If any BRI team is on Bitbucket Cloud rather than Bamboo, `bitbucket-pipelines.yml` is a much
closer analogue to Actions: container-per-step, `caches:`, `services:`, `parallel:`,
`deployment:`, and `pipes:` as the action equivalent. Bamboo is the older, agent-and-
capability-centric model, which is exactly why the on-prem/regulated setups still use it —
the permission and audit story is stronger.

---

Sources: [Bamboo Specs Reference](https://docs.atlassian.com/bamboo-specs-docs/9.0.0/specs.html?yaml)
and [Configuring a variables task](https://confluence.atlassian.com/spaces/BAMBOO058/pages/744331605/Configuring+a+variables+task).
Content was rephrased for compliance with licensing restrictions.
