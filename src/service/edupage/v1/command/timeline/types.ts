import { z } from 'zod';

import {
  edupageV1TimelineInputParameterSchema,
  edupageV1TimelineInputPayloadSchema,
  edupageV1TimelineOutputPayloadSchema,
} from './schema';

export type EdupageV1TimelineInputParameter = z.input<
  typeof edupageV1TimelineInputParameterSchema
>;

export type EdupageV1TimelineInputPayload = z.input<
  typeof edupageV1TimelineInputPayloadSchema
>;

export type EdupageV1TimelineOutputPayload = z.output<
  typeof edupageV1TimelineOutputPayloadSchema
>;
