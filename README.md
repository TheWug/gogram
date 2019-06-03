# Gogram

Gogram is a library for interacting with the Telegram Bot API in Golang.

## Getting started

Follow these instructions to use `gogram` in your Golang project

### Prerequisites

You should have Golang 1.10 or newer installed. It may work with older versions, but no testing or compatibility checking has been done.

### Installing

Just import `github.com/thewug/gogram` and fetch it like any other golang library.

### Using

1. Declare a `TelegramBot.`
2. Write yourself some callback providers that implement the interfaces in interfaces.go, and hook them to the bot.
  a. `Messagable` for handling ordinary messages sent to pms, groups, or channels
  b. `Callbackable` for handling callbacks from inline keyboard buttons
  c. `InlineQueryable` for handling inline queries
  d. `Maintainer` for periodic tasks
3. Use a `MessageStateMachine` to handle bot commands.  It supports both mundane commands where one command equals one action, and more complicated flows that request user input.
  a. Write yourself some command providers that implement the `State` interface and bind them to the state machine.
4. Write yourself an `InitSettings` implementation. The `InitSettings.InitializeAll` function should, at minimum, set the bot's API key, but it may do other things, such as setting up database connections, etc.
5. Call `TelegramBot.MainLoop` and off you go.

## Contributing

Please be a reasonable person. If you notice an actual bug, or if you want to add support for existing bot API features, you are welcome to submit a pull request. You are welcome to submit pull requests for large architectural changes as well but do so at your own peril, I may decide not to use them.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/TheWug/reqtify/tags).

## Authors

@TheWug

See also the list of [contributors](https://github.com/TheWug/gogram/contributors) who participated in this project.

## License

Copyright 2019 TheWug

This project is licensed under [the MIT license.](https://opensource.org/licenses/MIT)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
