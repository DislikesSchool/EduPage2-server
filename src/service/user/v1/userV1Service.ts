import {
  CommandDefinitionList,
  SubscriptionDefinitionList,
} from '@purista/core';

import { infoCommandBuilder } from './command/info';
import { registerCommandBuilder } from './command/register';
import { validateCommandBuilder } from './command/validate';
import { userV1ServiceBuilder } from './userV1ServiceBuilder';

// bring service config definition, command definitions and subscription definitions together in the service
// add only definitions and no further service config here
// other service config should be done in ./userServiceBuilder.ts file

const commandDefinitions: CommandDefinitionList<any> = [
  registerCommandBuilder.getDefinition(),
  validateCommandBuilder.getDefinition(),
  infoCommandBuilder.getDefinition(),
];

const subscriptionDefinitions: SubscriptionDefinitionList<any> = [];

export const userV1Service = userV1ServiceBuilder
  .addCommandDefinition(...commandDefinitions)
  .addSubscriptionDefinition(...subscriptionDefinitions);
