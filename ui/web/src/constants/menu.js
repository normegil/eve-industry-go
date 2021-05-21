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
},
{
  id: "pages",
  icon: "iconsminds-digital-drawing",
  label: "menu.pages",
  to: "/user/login",
  roles: [UserRole.Admin],
  subs: [{
    icon: "simple-icon-user-following",
    label: "menu.login",
    to: "/user/login",
    newWindow: true
  }, {
    icon: "simple-icon-user-follow",
    label: "menu.register",
    to: "/user/register",
    newWindow: true
  }, {
    icon: "simple-icon-user-unfollow",
    label: "menu.forgot-password",
    to: "/user/forgot-password",
    newWindow: true
  },
  {
    icon: "simple-icon-user-following",
    label: "menu.reset-password",
    to: "/user/reset-password",
    newWindow: true
  }
  ]
},
];
export default data;
