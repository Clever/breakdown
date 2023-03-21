## Link to JIRA:
[Link to JIRA](url)

## Overview:
[insert description]

If you made any changes to swagger.yml:
- [ ] Update swagger.yml version
- [ ] Run "make generate"

## Testing

## Rollout

## Rollback (in the event of a problem)
- [ ] Rollback via dapple OR `ark rollback -e production breakdown`?
- [ ] Anything else?

## New Repo Setup
- [ ] Set up Slack notifications for this app for your team https://clever.atlassian.net/wiki/spaces/ENG/pages/888897571/GitHub+assignments
- [ ] Follow the instructions at go/private-deps
- [ ] Tune your alarms if you'd like anything other than our default recommendations https://clever.atlassian.net/wiki/spaces/~620990898/pages/904036784/Alarm+Best+Practices
- [ ] Adjust the cpu and max_mem values based on the needs of your application as the containers are killed if these values are exceeded
- [ ] If this app does not need to be canaried, you can disable that in the `deploy_config` section of the launch config
- [ ] If this app should be multi region in sso pods adjust the pod_config https://clever.atlassian.net/wiki/spaces/ENG/pages/370147335/launch.yml#launch.yml-PodConfig
- [ ] Run "make install_deps" and commit the go.sum file
- [ ] Delete this section from the PR template
