import Keycloak from 'keycloak-js';

const keycloakUrl = import.meta.env.VITE_KEYCLOAK_URL;
const keycloakRealm = import.meta.env.VITE_KEYCLOAK_REALM;
const keycloakClientId = import.meta.env.VITE_KEYCLOAK_CLIENT_ID;
if (!keycloakUrl) throw new Error('VITE_KEYCLOAK_URL is not set. Set it in your .env file or pass it at build time.');
if (!keycloakRealm) throw new Error('VITE_KEYCLOAK_REALM is not set. Set it in your .env file or pass it at build time.');
if (!keycloakClientId) throw new Error('VITE_KEYCLOAK_CLIENT_ID is not set. Set it in your .env file or pass it at build time.');

const keycloak = new Keycloak({
    url: keycloakUrl,
    realm: keycloakRealm,
    clientId: keycloakClientId,
});

let isInitialized = false;

export async function initKeycloak() {
    if (isInitialized) return keycloak;

    try {
        const authenticated = await keycloak.init({
            onLoad: 'login-required',
        });

        isInitialized = true;
        return authenticated ? keycloak : null;
    } catch (err) {
        console.error('Keycloak init failed:', err);
        return null;
    }
}

export default keycloak;