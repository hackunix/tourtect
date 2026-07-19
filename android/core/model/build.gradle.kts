plugins {
    alias(libs.plugins.android.library)
    alias(libs.plugins.kotlin.android)
}

android {
    namespace = "com.tourtect.core.model"
    compileSdk = 35

    defaultConfig {
        minSdk = 26
    }
}

dependencies {
    testImplementation(libs.junit)
}
