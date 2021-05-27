import { adminRoot } from "./config";
import { UserRole } from "../utils/auth.roles";

const data = [{
  id: "Manufacturing",
  icon: "iconsminds-factory",
  label: "menu.manufacturing",
  to: `${adminRoot}/manufacturing`,
  roles: [UserRole.Admin, UserRole.User],
  subs: [{
    icon: "iconsminds-files",
    label: "menu.myblueprints",
    to: `${adminRoot}/manufacturing/myblueprints`,
  },
  ]
},
{
  id: "research",
  icon: "iconsminds-chemical",
  label: "menu.research",
  to: `${adminRoot}/research`,
  roles: [UserRole.Admin, UserRole.User],
  subs: [{
    icon: "simple-icon-paper-plane",
    label: "menu.invention",
    to: `${adminRoot}/research/invention`,
  },
  ]
},
{
  id: "market",
  icon: "iconsminds-line-chart-1",
  label: "menu.market",
  roles: [UserRole.Admin],
  to: `${adminRoot}/market`
},
{
  id: "settings",
  icon: "iconsminds-gear",
  label: "menu.settings",
  to: `${adminRoot}/settings`
}
];
export default data;
