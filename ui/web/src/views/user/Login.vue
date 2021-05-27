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
export default {};
</script>
