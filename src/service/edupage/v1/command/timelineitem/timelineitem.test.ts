import {
  getCommandContextMock,
  getEventBridgeMock,
  getLoggerMock,
} from '@purista/core';
import { createSandbox } from 'sinon';

import { edupageV1Service } from '../../edupageV1Service';
import { timelineitemCommandBuilder } from './timelineitemCommandBuilder';
import {
  EdupageV1TimelineitemInputParameter,
  EdupageV1TimelineitemInputPayload,
} from './types';

describe('service Edupage version 1 - command timelineitem', () => {
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

    const timelineitem = timelineitemCommandBuilder
      .getCommandFunction()
      .bind(service);

    const payload: EdupageV1TimelineitemInputPayload = undefined;

    const parameter: EdupageV1TimelineitemInputParameter = {};

    const context = getCommandContextMock(payload, parameter, sandbox);

    const result = await timelineitem(context.mock, payload, parameter);

    expect(result).toBeUndefined();
  });
});
