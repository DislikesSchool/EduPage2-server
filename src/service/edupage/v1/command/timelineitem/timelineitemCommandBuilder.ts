import { edupageV1ServiceBuilder } from '../../edupageV1ServiceBuilder';
import {
  edupageV1TimelineitemInputParameterSchema,
  edupageV1TimelineitemInputPayloadSchema,
  edupageV1TimelineitemOutputPayloadSchema,
} from './schema';

export const timelineitemCommandBuilder = edupageV1ServiceBuilder
  .getCommandBuilder(
    'timelineitem',
    'Get all the details of an item from the timeline by it&#x27;s ID',
  )
  .addPayloadSchema(edupageV1TimelineitemInputPayloadSchema)
  .addParameterSchema(edupageV1TimelineitemInputParameterSchema)
  .addOutputSchema(edupageV1TimelineitemOutputPayloadSchema)
  .exposeAsHttpEndpoint('GET', 'edupage/timeline/item')
  .addOpenApiTags('EduPage')
  .setCommandFunction(async function (_context, _payload, _parameter) {
    // add your business logic here
  });
