import { edupageV1ServiceBuilder } from '../../edupageV1ServiceBuilder';
import {
  edupageV1BasetimelineInputParameterSchema,
  edupageV1BasetimelineInputPayloadSchema,
  edupageV1BasetimelineOutputPayloadSchema,
} from './schema';

export const basetimelineCommandBuilder = edupageV1ServiceBuilder
  .getCommandBuilder(
    'basetimeline',
    'Returns all the timeline items, with no modification',
  )
  .addPayloadSchema(edupageV1BasetimelineInputPayloadSchema)
  .addParameterSchema(edupageV1BasetimelineInputParameterSchema)
  .addOutputSchema(edupageV1BasetimelineOutputPayloadSchema)
  .exposeAsHttpEndpoint('GET', 'edupage/timeline/base')
  .addOpenApiTags('EduPage')
  .setCommandFunction(async function (_context, _payload, _parameter) {
    // add your business logic here
  });
