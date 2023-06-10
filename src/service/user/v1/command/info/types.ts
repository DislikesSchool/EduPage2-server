import { z } from 'zod';

import {
  userV1InfoInputParameterSchema,
  userV1InfoInputPayloadSchema,
  userV1InfoOutputPayloadSchema,
} from './schema';

export type UserV1InfoInputParameter = z.input<
  typeof userV1InfoInputParameterSchema
>;

export type UserV1InfoInputPayload = z.input<
  typeof userV1InfoInputPayloadSchema
>;

export type UserV1InfoOutputPayload = z.output<
  typeof userV1InfoOutputPayloadSchema
>;
