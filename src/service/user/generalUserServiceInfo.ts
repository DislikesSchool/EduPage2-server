import { ServiceInfoType } from '@purista/core';

export const generalUserServiceInfo: Omit<ServiceInfoType, 'serviceVersion'> = {
  serviceName: 'User',
  serviceDescription:
    'Managing the user&#x27;s connection to their EduPage account',
};
