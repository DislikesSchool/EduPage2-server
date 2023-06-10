import { z } from 'zod';

import {
  edupageV1BasetimelineInputParameterSchema,
  edupageV1BasetimelineInputPayloadSchema,
  edupageV1BasetimelineOutputPayloadSchema,
} from './schema';

export type EdupageV1BasetimelineInputParameter = z.input<
  typeof edupageV1BasetimelineInputParameterSchema
>;

export type EdupageV1BasetimelineInputPayload = z.input<
  typeof edupageV1BasetimelineInputPayloadSchema
>;

export type EdupageV1BasetimelineOutputPayload = z.output<
  typeof edupageV1BasetimelineOutputPayloadSchema
>;
