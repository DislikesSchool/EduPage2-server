import { extendApi } from '@purista/core';
import { z } from 'zod';

// define the input parameters
export const icanteenV1SetupInputParameterSchema = extendApi(z.object({}), {
  title: 'setup input parameter schema',
});

// define the input payload
export const icanteenV1SetupInputPayloadSchema = extendApi(z.any(), {
  title: 'setup input payload schema',
});

// define the output payload
export const icanteenV1SetupOutputPayloadSchema = extendApi(z.any(), {
  title: 'setup output payload schema',
});
