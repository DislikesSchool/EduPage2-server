import { extendApi } from '@purista/core';
import { z } from 'zod';

// define the input parameters
export const userV1ValidateInputParameterSchema = extendApi(z.object({}), {
  title: 'validate input parameter schema',
});

// define the input payload
export const userV1ValidateInputPayloadSchema = extendApi(
  z.object({
    token: z.string().nonempty(),
  }),
  {
    title: 'validate input payload schema',
  },
);

// define the output payload
export const userV1ValidateOutputPayloadSchema = extendApi(
  z.object({
    valid: z.boolean(),
  }),
  {
    title: 'validate output payload schema',
  },
);
