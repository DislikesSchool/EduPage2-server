import { z } from 'zod';

import {
  edupageV1TimelineitemInputParameterSchema,
  edupageV1TimelineitemInputPayloadSchema,
  edupageV1TimelineitemOutputPayloadSchema,
} from './schema';

export type EdupageV1TimelineitemInputParameter = z.input<
  typeof edupageV1TimelineitemInputParameterSchema
>;

export type EdupageV1TimelineitemInputPayload = z.input<
  typeof edupageV1TimelineitemInputPayloadSchema
>;

export type EdupageV1TimelineitemOutputPayload = z.output<
  typeof edupageV1TimelineitemOutputPayloadSchema
>;
