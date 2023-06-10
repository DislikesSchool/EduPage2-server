import { ServiceBuilder, ServiceInfoType } from '@purista/core';

import { generalIcanteenServiceInfo } from '../generalIcanteenServiceInfo';
import { icanteenServiceV1ConfigSchema } from './icanteenServiceConfig';

export const icanteenServiceInfo: ServiceInfoType = {
  serviceVersion: '1',
  ...generalIcanteenServiceInfo,
};

// create a service builder instance and assign service config schema and default config.

export const icanteenV1ServiceBuilder = new ServiceBuilder(icanteenServiceInfo)
  .setConfigSchema(icanteenServiceV1ConfigSchema)
  .setDefaultConfig({});
