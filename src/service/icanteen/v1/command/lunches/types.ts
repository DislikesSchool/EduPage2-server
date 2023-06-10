import { z } from 'zod';

import {
  icanteenV1LunchesInputParameterSchema,
  icanteenV1LunchesInputPayloadSchema,
  icanteenV1LunchesOutputPayloadSchema,
} from './schema';

export type IcanteenV1LunchesInputParameter = z.input<
  typeof icanteenV1LunchesInputParameterSchema
>;

export type IcanteenV1LunchesInputPayload = z.input<
  typeof icanteenV1LunchesInputPayloadSchema
>;

export type IcanteenV1LunchesOutputPayload = z.output<
  typeof icanteenV1LunchesOutputPayloadSchema
>;
