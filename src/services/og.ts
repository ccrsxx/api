import { join } from 'path';
import satori from 'satori';
import { z } from 'zod';
import { html } from 'satori-html';
import { Resvg } from '@resvg/resvg-js';
import { getArrayBufferFromFile } from '../lib/helper.ts';
import { validStringSchema } from '../lib/validation.ts';

const validOgQuery = z.object({
  type: validStringSchema.optional(),
  title: validStringSchema.optional(),
  image: validStringSchema.optional(),
  article: validStringSchema.optional(),
  description: validStringSchema.optional()
});

type ValidOgQuery = z.infer<typeof validOgQuery>;

export async function getOg(query: ValidOgQuery): Promise<Buffer> {
  const parsedQuery = validOgQuery.safeParse(query);

  const { type, title, image, article, description } = parsedQuery.data ?? {};

  const isHomepage = title === 'Risal Amin';

  const markup = html(`
    <div tw='flex h-full w-full bg-black p-8 text-white'>
     ${
       article
         ? `
            <div tw='flex w-full justify-between'>
              <div tw='flex flex-col justify-between'>
                <div tw='flex flex-col'>
                  <p tw='-my-2 text-xl font-medium text-gray-400'>
                    risalamin.com/${type}
                  </p>
                  <h2 style=${gradientTitleStyles} tw='max-w-xl text-4xl'>
                    ${title}
                  </h2>
                </div>
                <div tw='flex items-center'>
                  <img
                    style='objectFit: cover;'
                    tw='h-18 w-18 rounded-full'
                    src='https://risalamin.com/assets/emilia.png'
                    alt='Emilia'
                  />
                  <div tw='ml-4 flex flex-col'>
                    <p tw='-mb-4 text-2xl font-semibold'>Risal Amin</p>
                    <p tw='text-lg font-medium text-gray-400'>@ccrsxx</p>
                  </div>
                </div>
              </div>
              <img
                style={{ objectFit: 'cover' }}
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
                src='https://risalamin.com/logo512.png'
                alt="risalamin.com's logo"
              />
              <h2 style=${gradientTitleStyles} tw='pb-1 text-6xl'>
                ${isHomepage ? 'Risal Amin' : title}
              </h2>
              ${
                !isHomepage
                  ? "<p tw='text-2xl font-semibold text-gray-200'>risalamin.com</p>"
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

  const svg = await satori(markup, {
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

  const resvg = new Resvg(svg);

  const pngData = resvg.render();
  const pngBuffer = pngData.asPng();

  return pngBuffer;
}

const gradientTitleStyles = `style="color: transparent; -webkit-background-clip: text; background-clip: text; 
                             background-image: linear-gradient(to right, #a855f7, #f472b6);"`;

const [interRegular, interMedium, interSemibold] = await Promise.all([
  getArrayBufferFromFile(
    join(import.meta.dirname, '../assets/inter-regular.ttf')
  ),
  getArrayBufferFromFile(
    join(import.meta.dirname, '../assets/inter-medium.ttf')
  ),
  getArrayBufferFromFile(
    join(import.meta.dirname, '../assets/inter-semibold.ttf')
  )
]);

export const OgService = {
  getOg
};
