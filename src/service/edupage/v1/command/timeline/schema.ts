import { extendApi } from '@purista/core';
import { z } from 'zod';

// define the input parameters
export const edupageV1TimelineInputParameterSchema = extendApi(z.object({}), {
  title: 'timeline input parameter schema',
});

// define the input payload
export const edupageV1TimelineInputPayloadSchema = extendApi(z.any(), {
  title: 'timeline input payload schema',
});

// define the output payload
export const edupageV1TimelineOutputPayloadSchema = extendApi(z.any(), {
  title: 'timeline output payload schema',
});
