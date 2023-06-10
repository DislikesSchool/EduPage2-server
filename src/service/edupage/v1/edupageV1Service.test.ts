import { edupageV1Service as service } from './edupageV1Service';

describe('service edupage version 1', () => {
  it('has valid commands', () => {
    service.validateCommandDefinitions();
  });

  it('has valid subscriptions', () => {
    service.validateSubscriptionDefinitions();
  });
});
