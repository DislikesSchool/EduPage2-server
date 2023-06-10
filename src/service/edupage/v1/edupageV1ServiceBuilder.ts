import { ServiceBuilder, ServiceInfoType } from '@purista/core';

import { generalEdupageServiceInfo } from '../generalEdupageServiceInfo';
import { edupageServiceV1ConfigSchema } from './edupageServiceConfig';

export const edupageServiceInfo: ServiceInfoType = {
  serviceVersion: '1',
  ...generalEdupageServiceInfo,
};

// create a service builder instance and assign service config schema and default config.

export const edupageV1ServiceBuilder = new ServiceBuilder(edupageServiceInfo)
  .setConfigSchema(edupageServiceV1ConfigSchema)
  .setDefaultConfig({});
