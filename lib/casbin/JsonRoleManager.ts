/* eslint-disable no-param-reassign */
/* eslint-disable no-restricted-syntax */
/* eslint-disable class-methods-use-this */
/* eslint-disable max-classes-per-file */
import { RoleManager, getLogger, logPrint } from 'casbin';
import _ from 'lodash';

// DEFAULT_DOMAIN defines the default domain space.
const DEFAULT_DOMAIN = 'casbin::default';

type MatchingFunc = (arg1: string, arg2: string) => boolean;

const getCombinations = (chars: string[]) => {
  const result: string[] = [];

  const f = function (prefix: string, remChars: string[]) {
    // eslint-disable-next-line no-plusplus
    for (let i = 0; i < remChars.length; i++) {
      const separator = prefix && remChars[i] ? '::' : '';
      const strToPush = prefix + separator + remChars[i];
      result.push(strToPush);
      f(strToPush, remChars.slice(i + 1));
    }
  };
  f('', chars);

  return result.map((v) => v.split('::'));
};

const getStringifiedJSONSubsets: (str: string) => string[] = (str: string) => {
  try {
    const obj = JSON.parse(str);
    const keyList: string[] = Object.keys(obj);
    const combinations = getCombinations(keyList);
    const stringifiedJsonSubsets = combinations.reduce((acc, keys) => {
      const subObj = _.pick(obj, keys);
      return [...acc, JSON.stringify(subObj)];
    }, []);
    return stringifiedJsonSubsets;
    // eslint-disable-next-line no-empty
  } catch (e) {}
  return [];
};

// loadOrDefault returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
function loadOrDefault<K, V>(map: Map<K, V>, key: K, value: V): V {
  const read = map.get(key);
  if (read === undefined) {
    map.set(key, value);
    return value;
  }
  return read;
}

/**
 * Role represents the data structure for a role in RBAC.
 */
class Role {
  public name: string;

  private roles: Role[];

  constructor(name: string) {
    this.name = name;
    this.roles = [];
  }

  public addRole(role: Role): void {
    if (this.roles.some((n) => n.name === role.name)) {
      return;
    }
    this.roles.push(role);
  }

  public deleteRole(role: Role): void {
    this.roles = this.roles.filter((n) => n.name !== role.name);
  }

  public hasRole(name: string, hierarchyLevel: number): boolean {
    if (this.name === name) {
      return true;
    }
    if (hierarchyLevel <= 0) {
      return false;
    }
    for (const role of this.roles) {
      if (role.hasRole(name, hierarchyLevel - 1)) {
        return true;
      }
    }

    return false;
  }

  public hasDirectRole(name: string): boolean {
    return this.roles.some((n) => n.name === name);
  }

  public toString(): string {
    return this.name + this.roles.join(', ');
  }

  public getRoles(): string[] {
    return this.roles.map((n) => n.name);
  }
}

class Roles extends Map<string, Role> {
  public hasRole(name: string, matchingFunc?: MatchingFunc): boolean {
    let ok = false;
    if (matchingFunc) {
      this.forEach((value, key) => {
        if (matchingFunc(name, key)) {
          ok = true;
        }
      });
    } else {
      return this.has(name);
    }
    return ok;
  }

  public createRole(name: string, matchingFunc?: MatchingFunc): Role {
    const role = loadOrDefault(this, name, new Role(name));
    if (matchingFunc) {
      this.forEach((value, key) => {
        if (matchingFunc(name, key) && name !== key) {
          // Add new role to matching role
          const role1 = loadOrDefault(this, key, new Role(key));
          role.addRole(role1);
        }
      });
    }
    return role;
  }
}

export class AttributeBasedRoleManager implements RoleManager {
  private allDomains: Map<string, Roles>;

  private maxHierarchyLevel: number;

  private hasPattern = false;

  private hasDomainPattern = false;

  private matchingFunc: MatchingFunc | undefined;

  private domainMatchingFunc: MatchingFunc | undefined;

  /**
   * DefaultRoleManager is the constructor for creating an instance of the
   * default RoleManager implementation.
   *
   * @param maxHierarchyLevel the maximized allowed RBAC hierarchy level.
   */
  constructor(maxHierarchyLevel: number) {
    this.allDomains = new Map<string, Roles>();
    this.allDomains.set(DEFAULT_DOMAIN, new Roles());
    this.maxHierarchyLevel = maxHierarchyLevel;
  }

  /**
   * addMatchingFunc support use pattern in g
   * @param name name
   * @param fn matching function
   * @deprecated
   */
  public async addMatchingFunc(name: string, fn: MatchingFunc): Promise<void>;

  /**
   * addMatchingFunc support use pattern in g
   * @param fn matching function
   */
  public async addMatchingFunc(fn: MatchingFunc): Promise<void>;

  /**
   * addMatchingFunc support use pattern in g
   * @param name name
   * @param fn matching function
   * @deprecated
   */
  public async addMatchingFunc(
    name: string | MatchingFunc,
    fn?: MatchingFunc
  ): Promise<void> {
    this.hasPattern = true;
    if (typeof name === 'string' && fn) {
      this.matchingFunc = fn;
    } else if (typeof name === 'function') {
      this.matchingFunc = name;
    } else {
      throw new Error('error: domain should be 1 parameter');
    }
  }

  /**
   * addDomainMatchingFunc support use domain pattern in g
   * @param fn domain matching function
   * ```
   */
  public async addDomainMatchingFunc(fn: MatchingFunc): Promise<void> {
    this.hasDomainPattern = true;
    this.domainMatchingFunc = fn;
  }

  private generateTempRoles(domain: string): Roles {
    loadOrDefault(this.allDomains, domain, new Roles());

    const patternDomain = new Set([domain]);
    if (this.hasDomainPattern) {
      this.allDomains.forEach((value, key) => {
        if (this.domainMatchingFunc && this.domainMatchingFunc(domain, key)) {
          patternDomain.add(key);
        }
      });
    }

    const allRoles = new Roles();
    patternDomain.forEach((d) => {
      loadOrDefault(this.allDomains, d, new Roles()).forEach((value) => {
        const role1 = allRoles.createRole(value.name, this.matchingFunc);
        value.getRoles().forEach((n) => {
          role1.addRole(allRoles.createRole(n, this.matchingFunc));
        });
      });
    });
    return allRoles;
  }

  /**
   * addLink adds the inheritance link between role: name1 and role: name2.
   * aka role: name1 inherits role: name2.
   * domain is a prefix to the roles.
   */
  public async addLink(
    name1: string,
    name2: string,
    ...domain: string[]
  ): Promise<void> {
    if (domain.length === 0) {
      domain = [DEFAULT_DOMAIN];
    } else if (domain.length > 1) {
      throw new Error('error: domain should be 1 parameter');
    }

    const allRoles = loadOrDefault(this.allDomains, domain[0], new Roles());

    const role1 = loadOrDefault(allRoles, name1, new Role(name1));
    const role2 = loadOrDefault(allRoles, name2, new Role(name2));
    role1.addRole(role2);
  }

  /**
   * clear clears all stored data and resets the role manager to the initial state.
   */
  public async clear(): Promise<void> {
    this.allDomains = new Map();
    this.allDomains.set(DEFAULT_DOMAIN, new Roles());
  }

  /**
   * deleteLink deletes the inheritance link between role: name1 and role: name2.
   * aka role: name1 does not inherit role: name2 any more.
   * domain is a prefix to the roles.
   */
  public async deleteLink(
    name1: string,
    name2: string,
    ...domain: string[]
  ): Promise<void> {
    if (domain.length === 0) {
      domain = [DEFAULT_DOMAIN];
    } else if (domain.length > 1) {
      throw new Error('error: domain should be 1 parameter');
    }

    const allRoles = loadOrDefault(this.allDomains, domain[0], new Roles());

    if (!allRoles.has(name1) || !allRoles.has(name2)) {
      return;
    }

    const role1 = loadOrDefault(allRoles, name1, new Role(name1));
    const role2 = loadOrDefault(allRoles, name2, new Role(name2));
    role1.deleteRole(role2);
  }

  private addJSONSubsetsAsRoles(allRoles: Roles, roleName: string) {
    // we need to dynamically add name1 and name2 so that we can compare attributes that are not present in the policy
    if (!allRoles.get(roleName)) {
      allRoles.set(roleName, new Role(roleName));
    }

    const roleToAddSubsetsTo = allRoles.get(roleName);

    const getAllRelatedRoleNamesInHierarchy: (
      role: Role,
      result?: Set<string>
    ) => Set<string> = (role: Role, result: Set<string> = new Set()) => {
      if (result.has(role.name)) return result;
      result.add(role.name);

      const immediateRoleNames = role.getRoles();
      immediateRoleNames.forEach((immediateRoleName) => {
        const immediateRole = allRoles.get(immediateRoleName);
        if (immediateRole) {
          // eslint-disable-next-line @typescript-eslint/no-unused-vars
          getAllRelatedRoleNamesInHierarchy(immediateRole, result);
        }
      });

      return result;
    };

    if (roleToAddSubsetsTo) {
      const relatedRolesNames = getAllRelatedRoleNamesInHierarchy(
        roleToAddSubsetsTo
      );
      relatedRolesNames.forEach((relatedRoleName) => {
        const relatedRole = allRoles.get(relatedRoleName);
        if (relatedRole) {
          const stringifiedJsonSubsets = getStringifiedJSONSubsets(
            relatedRoleName
          );
          stringifiedJsonSubsets.forEach((jsonSubsetStr) => {
            const jsonSubsetRole =
              allRoles.get(jsonSubsetStr) || new Role(jsonSubsetStr);
            relatedRole.addRole(jsonSubsetRole);
            allRoles.set(jsonSubsetStr, jsonSubsetRole);
          });
        }
      });
    }
    return allRoles;
  }

  /**
   * hasLink determines whether role: name1 inherits role: name2.
   * domain is a prefix to the roles.
   */
  public async hasLink(
    name1: string,
    name2: string,
    ...domain: string[]
  ): Promise<boolean> {
    if (domain.length === 0) {
      // eslint-disable-next-line no-param-reassign
      domain = [DEFAULT_DOMAIN];
    } else if (domain.length > 1) {
      throw new Error('error: domain should be 1 parameter');
    }

    if (name1 === name2) {
      return true;
    }

    let allRoles: Roles;
    if (this.hasPattern || this.hasDomainPattern) {
      allRoles = this.generateTempRoles(domain[0]);
    } else {
      allRoles = loadOrDefault(this.allDomains, domain[0], new Roles());
    }

    allRoles = this.addJSONSubsetsAsRoles(allRoles, name1);

    if (
      !allRoles.hasRole(name1, this.matchingFunc) ||
      !allRoles.hasRole(name2, this.matchingFunc)
    ) {
      return false;
    }

    const role1 = allRoles.createRole(name1, this.matchingFunc);
    return role1.hasRole(name2, this.maxHierarchyLevel);
  }

  /**
   * getRoles gets the roles that a subject inherits.
   * domain is a prefix to the roles.
   */
  public async getRoles(name: string, ...domain: string[]): Promise<string[]> {
    if (domain.length === 0) {
      domain = [DEFAULT_DOMAIN];
    } else if (domain.length > 1) {
      throw new Error('error: domain should be 1 parameter');
    }

    let allRoles: Roles;
    if (this.hasPattern || this.hasDomainPattern) {
      allRoles = this.generateTempRoles(domain[0]);
    } else {
      allRoles = loadOrDefault(this.allDomains, domain[0], new Roles());
    }

    if (!allRoles.hasRole(name, this.matchingFunc)) {
      return [];
    }

    return allRoles.createRole(name, this.matchingFunc).getRoles();
  }

  /**
   * getUsers gets the users that inherits a subject.
   * domain is an unreferenced parameter here, may be used in other implementations.
   */
  public async getUsers(name: string, ...domain: string[]): Promise<string[]> {
    if (domain.length === 0) {
      domain = [DEFAULT_DOMAIN];
    } else if (domain.length > 1) {
      throw new Error('error: domain should be 1 parameter');
    }

    let allRoles: Roles;
    if (this.hasPattern || this.hasDomainPattern) {
      allRoles = this.generateTempRoles(domain[0]);
    } else {
      allRoles = loadOrDefault(this.allDomains, domain[0], new Roles());
    }

    if (!allRoles.hasRole(name, this.matchingFunc)) {
      return [];
    }

    return [...allRoles.values()]
      .filter((n) => n.hasDirectRole(name))
      .map((n) => n.name);
  }

  /**
   * printRoles prints all the roles to log.
   */
  public async printRoles(): Promise<void> {
    if (getLogger().isEnable()) {
      [...this.allDomains.values()].forEach((n) => {
        logPrint(n.toString());
      });
    }
  }
}