import {
  getCommandContextMock,
  getEventBridgeMock,
  getLoggerMock,
} from '@purista/core';
import { createSandbox } from 'sinon';

import { edupageV1Service } from '../../edupageV1Service';
import { timelineCommandBuilder } from './timelineCommandBuilder';
import {
  EdupageV1TimelineInputParameter,
  EdupageV1TimelineInputPayload,
} from './types';

describe('service Edupage version 1 - command timeline', () => {
  let sandbox = createSandbox();
  beforeEach(() => {
    sandbox = createSandbox();
  });

  afterEach(() => {
    sandbox.restore();
  });

  test('does not throw', async () => {
    const service = edupageV1Service.getInstance(
      getEventBridgeMock(sandbox).mock,
      { logger: getLoggerMock(sandbox).mock },
    );

    const timeline = timelineCommandBuilder.getCommandFunction().bind(service);

    const payload: EdupageV1TimelineInputPayload = undefined;

    const parameter: EdupageV1TimelineInputParameter = {};

    const context = getCommandContextMock(payload, parameter, sandbox);

    const result = await timeline(context.mock, payload, parameter);

    expect(result).toBeUndefined();
  });
});
