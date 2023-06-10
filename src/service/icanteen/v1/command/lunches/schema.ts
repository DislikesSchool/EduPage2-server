import { extendApi } from '@purista/core';
import { z } from 'zod';

// define the input parameters
export const icanteenV1LunchesInputParameterSchema = extendApi(z.object({}), {
  title: 'lunches input parameter schema',
});

// define the input payload
export const icanteenV1LunchesInputPayloadSchema = extendApi(z.any(), {
  title: 'lunches input payload schema',
});

// define the output payload
export const icanteenV1LunchesOutputPayloadSchema = extendApi(z.any(), {
  title: 'lunches output payload schema',
});
