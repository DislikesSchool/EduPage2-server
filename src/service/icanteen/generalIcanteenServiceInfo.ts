import { ServiceInfoType } from '@purista/core';

export const generalIcanteenServiceInfo: Omit<
  ServiceInfoType,
  'serviceVersion'
> = {
  serviceName: 'Icanteen',
  serviceDescription: 'Manages iCanteen integration for users',
};
