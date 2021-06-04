import {UserRole} from "../utils/auth.roles";

export const defaultMenuType = 'menu-default' // 'menu-default', 'menu-sub-hidden', 'menu-hidden';
export const adminRoot = '';
export const searchPath = `${adminRoot}/#`
export const buyUrl = 'https://1.envato.market/nEyZa'

export function apiUrl() {
    let url = process.env.VUE_APP_API_BASE_URL
    if (url === undefined || url === null || url === "") {
        return window.location.origin
    }
    return url
}

export const subHiddenBreakpoint = 1440
export const menuHiddenBreakpoint = 768

export const defaultLocale = 'en'
export const defaultDirection = 'ltr'
export const localeOptions = [
    {id: 'en', name: 'English LTR', direction: 'ltr'}
]

export const firebaseConfig = {
    apiKey: "AIzaSyDe94G7L-3soXVSpVbsYlB5DfYH2L91aTU",
    authDomain: "piaf-vue-login.firebaseapp.com",
    databaseURL: "https://piaf-vue-login.firebaseio.com",
    projectId: "piaf-vue-login",
    storageBucket: "piaf-vue-login.appspot.com",
    messagingSenderId: "557576321564",
    appId: "1:557576321564:web:bc2ce73477aff5c2197dd9"
};


export const currentUser = {
    id: -1,
    name: 'Anonymous',
    portraits: {
        url64: '/assets/img/profiles/l-1.jpg',
        url128: '/assets/img/profiles/l-1.jpg',
        url256: '/assets/img/profiles/l-1.jpg',
        url512: '/assets/img/profiles/l-1.jpg'
    },
    date: 'Last seen today 15:24',
    role: UserRole.None
}

export const isAuthGuardActive = true;
export const themeRadiusStorageKey = 'theme_radius'
export const themeSelectedColorStorageKey = 'theme_selected_color'
export const defaultColor = 'light.blueolympic'
export const colors = ['bluenavy', 'blueyale', 'blueolympic', 'greenmoss', 'greenlime', 'purplemonster', 'orangecarrot', 'redruby', 'yellowgranola', 'greysteel']
