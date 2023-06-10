export enum ServiceEvent {
  /**
   * Emitted by user v1 command register:
   * Registers a new user by saving their EduPage credentials
   */
  UserRegistered = 'userRegistered',
  /**
   * Emitted by icanteen v1 command setup:
   * Sets up the iCanteen integration for a user
   */
  IcanteenSetup = 'icanteenSetup',
}
