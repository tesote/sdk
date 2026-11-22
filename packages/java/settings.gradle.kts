plugins {
    id("com.gradleup.nmcp.settings") version "1.6.2"
}

rootProject.name = "sdk"

nmcpSettings {
    centralPortal {
        username.set(System.getenv("MAVEN_CENTRAL_USERNAME") ?: "")
        password.set(System.getenv("MAVEN_CENTRAL_PASSWORD") ?: "")
        publishingType.set("AUTOMATIC")
    }
}
