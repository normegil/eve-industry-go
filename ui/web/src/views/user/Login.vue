<template>
<b-row class="h-100">
    <b-colxx xxs="12" md="10" class="mx-auto my-auto">
        <b-card class="auth-card" no-body>
            <div class="position-relative image-side">
                <p class="text-white h2">{{ $t('dashboards.magic-is-in-the-details') }}</p>
            </div>
            <div class="form-side">
                <router-link tag="a" to="/">
                    <span class="logo-single" />
                </router-link>

                <a href="/auth/login">
                  <img src="/assets/img/login/eve-sso-login-white-large.png" />
                </a>
            </div>
        </b-card>
    </b-colxx>
</b-row>
</template>

<script>
import {
    mapGetters,
    mapActions
} from "vuex";
import {
    validationMixin
} from "vuelidate";
const {
    required,
    maxLength,
    minLength,
    email
} = require("vuelidate/lib/validators");
import { adminRoot } from '../../constants/config';

export default {
    data() {
        return {
            form: {
                email: "test@coloredstrategies.com",
                password: "xxxxxx"
            },
        };
    },
    mixins: [validationMixin],
    validations: {
        form: {
            password: {
                required,
                maxLength: maxLength(16),
                minLength: minLength(4)
            },
            email: {
                required,
                email,
                minLength: minLength(4)
            }
        }
    },
    computed: {
        ...mapGetters(["currentUser", "processing", "loginError"])
    },
    methods: {
        ...mapActions(["login"]),
        formSubmit() {
            this.$v.$touch();
            this.form.email = "piaf-vue@coloredstrategies.com";
            this.form.password = "piaf123";
            this.$v.form.$touch();
           // if (!this.$v.form.$anyError) {
                this.login({
                    email: this.form.email,
                    password: this.form.password
                });
            //}
        }
    },
    watch: {
        currentUser(val) {
            if (val && val.uid && val.uid.length > 0) {
                setTimeout(() => {
                    this.$router.push(adminRoot);
                }, 200);
            }
        },
        loginError(val) {
            if (val != null) {
                this.$notify("error", "Login Error", val, {
                    duration: 3000,
                    permanent: false
                });

            }
        }
    }
};
</script>
