import { edupageV1ServiceBuilder } from '../../edupageV1ServiceBuilder';
import {
  edupageV1TimelineInputParameterSchema,
  edupageV1TimelineInputPayloadSchema,
  edupageV1TimelineOutputPayloadSchema,
} from './schema';

export const timelineCommandBuilder = edupageV1ServiceBuilder
  .getCommandBuilder(
    'timeline',
    'Returns optimised timeline, stripped of useless info',
  )
  .addPayloadSchema(edupageV1TimelineInputPayloadSchema)
  .addParameterSchema(edupageV1TimelineInputParameterSchema)
  .addOutputSchema(edupageV1TimelineOutputPayloadSchema)
  .exposeAsHttpEndpoint('GET', 'edupage/timeline')
  .addOpenApiTags('EduPage')
  .setCommandFunction(async function (_context, _payload, _parameter) {
    // add your business logic here
  });
