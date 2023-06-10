import { z } from 'zod';

// define the service config schema and the default service configuration

export const icanteenServiceV1ConfigSchema = z.object({});

export type IcanteenServiceV1Config = z.input<
  typeof icanteenServiceV1ConfigSchema
>;
