plugins {
    id("com.gradleup.nmcp.settings") version "1.6.2"
}

rootProject.name = "sdk"

nmcpSettings {
    centralPortal {
        username.set(providers.environmentVariable("MAVEN_CENTRAL_USERNAME"))
        password.set(providers.environmentVariable("MAVEN_CENTRAL_PASSWORD"))
        publishingType.set("AUTOMATIC")
    }
}
