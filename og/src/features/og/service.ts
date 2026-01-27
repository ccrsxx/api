import { z } from 'zod';
import satori from 'satori';
import { html } from 'satori-html';
import { Resvg } from '@resvg/resvg-js';
import { appEnv } from '../../config/env.ts';
import { AppError } from '../../api/error.ts';
import { validStringSchema } from '../../utils/validation.ts';
import type { ParsedQs } from 'qs';

const emiliaImageUrl = `${appEnv.CLOUDFLARE_CDN_URL}/assets/emilia.png`;

const validOgQuery = z.object({
  type: validStringSchema.optional(),
  title: validStringSchema.optional().default(''),
  image: z
    .string()
    .transform((val) => z.url().safeParse(val).data)
    .optional()
    .default(emiliaImageUrl),
  article: validStringSchema.optional(),
  description: validStringSchema.optional().default('')
});

export type ValidOgQuery = z.infer<typeof validOgQuery>;

async function getOg(query: ParsedQs): Promise<Buffer> {
  const { type, title, image, article, description } =
    validOgQuery.safeParse(query).data ?? {};

  const isHomepage = title === appEnv.OWNER_NAME;

  const markup = html(`
    <div tw='flex h-full w-full bg-black p-8 text-white'>
     ${
       article
         ? `
            <div tw='flex w-full justify-between'>
              <div tw='flex flex-col justify-between'>
                <div tw='flex flex-col'>
                  <p tw='-my-2 text-xl font-medium text-gray-400'>
                    ${appEnv.FRONTEND_DOMAIN}/${type}
                  </p>
                  <h2 style=${gradientTitleStyles} tw='max-w-xl text-4xl'>
                    ${title}
                  </h2>
                </div>
                <div tw='flex items-center'>
                  <img
                    style='object-fit: cover;'
                    tw='h-18 w-18 rounded-full'
                    src='${emiliaImageUrl}'
                    alt='Emilia'
                  />
                  <div tw='ml-4 flex flex-col'>
                    <p tw='-mb-4 text-2xl font-semibold'>${appEnv.OWNER_NAME}</p>
                    <p tw='text-lg font-medium text-gray-400'>@ccrsxx</p>
                  </div>
                </div>
              </div>
              <img
                style='object-fit: cover;'
                tw='h-full w-[448px] rounded-md'
                src='${image}'
                alt='${title}'
              />
            </div>
          `
         : `
            <div tw='flex w-full flex-col items-center justify-center'>
              <img
                tw='h-24 w-24'
                src='https://${appEnv.FRONTEND_DOMAIN}/logo512.png'
                alt="${appEnv.FRONTEND_DOMAIN}'s logo"
              />
              <h2 style=${gradientTitleStyles} tw='pb-1 text-6xl'>
                ${isHomepage ? 'Risal Amin' : title}
              </h2>
              ${
                !isHomepage
                  ? `<p tw='text-2xl font-semibold text-gray-200'>${appEnv.FRONTEND_DOMAIN}</p>`
                  : ''
              }
              <p tw='max-w-4xl text-center text-2xl text-gray-300'>
                ${description}
              </p>
            </div>
        `
     }
    </div>
  `);

  let svg: string;

  try {
    svg = await satori(markup, {
      width: 1200,
      height: 600,
      fonts: [
        {
          name: 'Inter',
          data: interRegular,
          weight: 400
        },
        {
          name: 'Inter',
          data: interMedium,
          weight: 500
        },
        {
          name: 'Inter',
          data: interSemibold,
          weight: 600
        }
      ]
    });
  } catch (err) {
    throw new AppError({
      message: 'Failed to generate OG image SVG.',
      cause: err
    });
  }

  const resvg = new Resvg(svg);

  const pngData = resvg.render();
  const pngBuffer = pngData.asPng();

  return pngBuffer;
}

const gradientTitleStyles = `style='color: transparent; -webkit-background-clip: text; background-clip: text; 
                             background-image: linear-gradient(to right, #a855f7, #f472b6);'`;

const [interRegular, interMedium, interSemibold] = await Promise.all([
  fetch(`${appEnv.CLOUDFLARE_CDN_URL}/assets/inter-regular.ttf`).then((res) =>
    res.arrayBuffer()
  ),
  fetch(`${appEnv.CLOUDFLARE_CDN_URL}/assets/inter-medium.ttf`).then((res) =>
    res.arrayBuffer()
  ),
  fetch(`${appEnv.CLOUDFLARE_CDN_URL}/assets/inter-semibold.ttf`).then((res) =>
    res.arrayBuffer()
  )
]);

export const OgService = {
  getOg
};
