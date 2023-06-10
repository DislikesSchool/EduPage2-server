import {
  CommandDefinitionList,
  SubscriptionDefinitionList,
} from '@purista/core';

import { basetimelineCommandBuilder } from './command/basetimeline';
import { timelineCommandBuilder } from './command/timeline';
import { timelineitemCommandBuilder } from './command/timelineitem';
import { edupageV1ServiceBuilder } from './edupageV1ServiceBuilder';

// bring service config definition, command definitions and subscription definitions together in the service
// add only definitions and no further service config here
// other service config should be done in ./edupageServiceBuilder.ts file

const commandDefinitions: CommandDefinitionList<any> = [
  basetimelineCommandBuilder.getDefinition(),
  timelineCommandBuilder.getDefinition(),
  timelineitemCommandBuilder.getDefinition(),
];

const subscriptionDefinitions: SubscriptionDefinitionList<any> = [];

export const edupageV1Service = edupageV1ServiceBuilder
  .addCommandDefinition(...commandDefinitions)
  .addSubscriptionDefinition(...subscriptionDefinitions);
