import { ServiceInfoType } from '@purista/core';

export const generalEdupageServiceInfo: Omit<
  ServiceInfoType,
  'serviceVersion'
> = {
  serviceName: 'Edupage',
  serviceDescription: 'Interfacing with the EduPage proxy and cache system',
};
