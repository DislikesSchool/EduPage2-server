import { z } from 'zod';

import {
  userV1ValidateInputParameterSchema,
  userV1ValidateInputPayloadSchema,
  userV1ValidateOutputPayloadSchema,
} from './schema';

export type UserV1ValidateInputParameter = z.input<
  typeof userV1ValidateInputParameterSchema
>;

export type UserV1ValidateInputPayload = z.input<
  typeof userV1ValidateInputPayloadSchema
>;

export type UserV1ValidateOutputPayload = z.output<
  typeof userV1ValidateOutputPayloadSchema
>;
