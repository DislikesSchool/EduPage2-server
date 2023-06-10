import {
  CommandDefinitionList,
  SubscriptionDefinitionList,
} from '@purista/core';

import { lunchesCommandBuilder } from './command/lunches';
import { setupCommandBuilder } from './command/setup';
import { icanteenV1ServiceBuilder } from './icanteenV1ServiceBuilder';

// bring service config definition, command definitions and subscription definitions together in the service
// add only definitions and no further service config here
// other service config should be done in ./icanteenServiceBuilder.ts file

const commandDefinitions: CommandDefinitionList<any> = [
  setupCommandBuilder.getDefinition(),
  setupCommandBuilder.getDefinition(),
  lunchesCommandBuilder.getDefinition(),
];

const subscriptionDefinitions: SubscriptionDefinitionList<any> = [];

export const icanteenV1Service = icanteenV1ServiceBuilder
  .addCommandDefinition(...commandDefinitions)
  .addSubscriptionDefinition(...subscriptionDefinitions);
