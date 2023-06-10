import { extendApi } from '@purista/core';
import { z } from 'zod';

// define the input parameters
export const edupageV1BasetimelineInputParameterSchema = extendApi(
  z.object({}),
  { title: 'basetimeline input parameter schema' },
);

// define the input payload
export const edupageV1BasetimelineInputPayloadSchema = extendApi(z.any(), {
  title: 'basetimeline input payload schema',
});

// define the output payload
export const edupageV1BasetimelineOutputPayloadSchema = extendApi(z.any(), {
  title: 'basetimeline output payload schema',
});
