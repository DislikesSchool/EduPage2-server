import { icanteenV1Service as service } from './icanteenV1Service';

describe('service icanteen version 1', () => {
  it('has valid commands', () => {
    service.validateCommandDefinitions();
  });

  it('has valid subscriptions', () => {
    service.validateSubscriptionDefinitions();
  });
});
