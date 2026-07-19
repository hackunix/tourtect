pluginManagement {
    repositories {
        google()
        mavenCentral()
        gradlePluginPortal()
    }
}
dependencyResolutionManagement {
    repositoriesMode.set(RepositoriesMode.FAIL_ON_PROJECT_REPOS)
    repositories {
        google()
        mavenCentral()
    }
}

rootProject.name = "Tourtect"
include(":app")
include(":core:network")
include(":core:database")
include(":core:security")
include(":core:designsystem")
include(":core:model")
include(":feature-forum")
include(":feature-assistant")
include(":feature-safety")
include(":feature-live")
include(":feature-lens")
