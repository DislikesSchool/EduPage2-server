import { extendApi } from '@purista/core';
import { z } from 'zod';

// define the input parameters
export const userV1InfoInputParameterSchema = extendApi(z.object({}), {
  title: 'info input parameter schema',
});

// define the input payload
export const userV1InfoInputPayloadSchema = extendApi(z.any(), {
  title: 'info input payload schema',
});

// define the output payload
export const userV1InfoOutputPayloadSchema = extendApi(z.any(), {
  title: 'info output payload schema',
});
