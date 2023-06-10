import { extendApi } from '@purista/core';
import { z } from 'zod';

// define the input parameters
export const edupageV1TimelineitemInputParameterSchema = extendApi(
  z.object({}),
  { title: 'timelineitem input parameter schema' },
);

// define the input payload
export const edupageV1TimelineitemInputPayloadSchema = extendApi(z.any(), {
  title: 'timelineitem input payload schema',
});

// define the output payload
export const edupageV1TimelineitemOutputPayloadSchema = extendApi(z.any(), {
  title: 'timelineitem output payload schema',
});
