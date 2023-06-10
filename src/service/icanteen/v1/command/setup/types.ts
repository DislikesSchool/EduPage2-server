import { z } from 'zod';

import {
  icanteenV1SetupInputParameterSchema,
  icanteenV1SetupInputPayloadSchema,
  icanteenV1SetupOutputPayloadSchema,
} from './schema';

export type IcanteenV1SetupInputParameter = z.input<
  typeof icanteenV1SetupInputParameterSchema
>;

export type IcanteenV1SetupInputPayload = z.input<
  typeof icanteenV1SetupInputPayloadSchema
>;

export type IcanteenV1SetupOutputPayload = z.output<
  typeof icanteenV1SetupOutputPayloadSchema
>;
