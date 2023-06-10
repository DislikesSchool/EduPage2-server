import { z } from 'zod';

// define the service config schema and the default service configuration

export const edupageServiceV1ConfigSchema = z.object({});

export type EdupageServiceV1Config = z.input<
  typeof edupageServiceV1ConfigSchema
>;
