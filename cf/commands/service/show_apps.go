package service

import (
        . "github.com/cloudfoundry/cli/cf/i18n"

        "github.com/cloudfoundry/cli/cf/api"
        "github.com/cloudfoundry/cli/cf/command_metadata"
        "github.com/cloudfoundry/cli/cf/configuration"
        "github.com/cloudfoundry/cli/cf/requirements"
        "github.com/cloudfoundry/cli/cf/terminal"
        "github.com/codegangsta/cli"
)

type ListAppsForService struct {
        ui                 terminal.UI
        config             configuration.Reader
        serviceSummaryRepo api.ServiceSummaryRepository
}

func NewListAppsForService(ui terminal.UI, config configuration.Reader, serviceSummaryRepo api.ServiceSummaryRepository) (cmd ListAppsForService) {
        cmd.ui = ui
        cmd.config = config
        cmd.serviceSummaryRepo = serviceSummaryRepo
        return
}

func (cmd ListAppsForService) Metadata() command_metadata.CommandMetadata {
        return command_metadata.CommandMetadata{
                Name:        "show-apps",
                Description: T("List all apps using same service instance in the target space"),
                Usage:       "CF_NAME show-apps SERVICE",
        }
}

func (cmd ListAppsForService) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
        reqs = append(reqs,
                requirementsFactory.NewLoginRequirement(),
                requirementsFactory.NewTargetedSpaceRequirement(),
        )
        return
}

func (cmd ListAppsForService) Run(c *cli.Context) {
        cmd.ui.Say(T("Getting services in org {{.OrgName}} / space {{.SpaceName}} as {{.CurrentUser}}...",
                map[string]interface{}{
                        "OrgName":     terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
                        "SpaceName":   terminal.EntityNameColor(cmd.config.SpaceFields().Name),
                        "CurrentUser": terminal.EntityNameColor(cmd.config.Username()),
                }))
                serviceInstances, apiErr := cmd.serviceSummaryRepo.GetSummariesInCurrentSpace()

        cmd.ui.Ok()
        cmd.ui.Say("")
        if apiErr != nil {
                cmd.ui.Failed(apiErr.Error())
                return
        }


        table := terminal.NewTable(cmd.ui, []string{T("Application names")})

        for _, instance := range serviceInstances {
                  serviceName := c.Args()[0]

                if instance.Name == serviceName{

                        for _, name := range instance.ApplicationNames{
                                table.Add(name)
                        }

                }
        }



        table.Print()
}

