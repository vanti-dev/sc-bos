/**
 * @returns {Promise<string | null>}
 */
export async function apiToken() {
  return localStorage.getItem("token");
}
