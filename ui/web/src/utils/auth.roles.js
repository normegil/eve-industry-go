export const UserRole = {
  None: -1,
  User: 1,

};

export function toRole(roleName) {
  if ("user" === roleName) {
    return UserRole.User
  }
  return UserRole.None;
}