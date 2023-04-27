# slack-openai

This is just me playing with [encore.dev](https://encore.dev/) and [OpenAI](https://openai.com/) to see if I can make a
chatbot that knows about my workplace.

## TODO

- Use encore.dev's database stuff to track conversations per Slack user so they can have the expected ChatGPT-style
  conversation (rather than single messages)
- Try to find a way to provide some sort of data about something particular to my workplace as context to the bot
    - e.g. "Hey, what's the IP address for the TeamCity server?"

## The good

- It's easy to spin up and get going
- The local dev environment saves you writing bespoke tooling to test your code locally
    - And it seems to work for all aspects of the encore.dev offering (APIs, Pub/Sub, Databases etc)
- The managed CI/CD saves you running and maintaining a CI system
- The secrets management is nice and easy to use (which will promote its usage)
- The "migrate away" enablement is good; if you got into trouble you could eject your apps as a containers and straight
  up run them in ECS on a dedicated EC2 host while you figured out how to decouple everything

## The bad

- It's prescriptive and opaque
- It needs pretty open access to your AWS account to do its thing
- It's pretty slow to respond to requests
- If you need to deviate from it's (limited) offering, it's not clear how you'd couple non-encore.dev resources to
  encore.dev resources
    - e.g. Maybe I want a single Lambda in front of it all that fires back an empty 200 for any request and I'm willing
      to pay to keep that Lambda hot
- Any change to the watched repo causes a deployment, so either you run a monorepo app to make it easy to reuse code and
  any minor change to a single endpoint causes a rollout to all endpoints, or you fragment your code into multiple
  repos / apps and have to manage the dependencies between them 
