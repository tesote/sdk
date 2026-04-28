plugins {
    `java-library`
    `maven-publish`
    signing
    // why: 0.1.0 has a broken transitive snapshot dep on
    // com.gradleup.gratatouille:gratatouille-runtime that's no longer in
    // the plugin portal; pinning to 0.0.9 (last good release).
    id("com.gradleup.nmcp") version "0.0.9"
}

group = "com.tesote"
version = "0.1.1"
description = "Official Java SDK for the equipo.tesote.com API"

java {
    toolchain {
        languageVersion.set(JavaLanguageVersion.of(17))
    }
    withJavadocJar()
    withSourcesJar()
}

repositories {
    mavenCentral()
}

dependencies {
    // why: stdlib java.net.http handles HTTP; jackson-databind is the single
    // runtime dep because jakarta.json's pull-style API is too awkward for
    // the dynamic envelopes the API returns. Justified in README.
    //
    // Version range (not `2.+`): Maven Central's component validation
    // rejects POMs that contain `+` in dependency versions. `[2.18,3)` is
    // the equivalent — any 2.x ≥ 2.18, never 3.x. Same shape used below for
    // test deps so all generated POMs are Central-compatible.
    api("com.fasterxml.jackson.core:jackson-databind:[2.18,3)")

    testImplementation("org.junit.jupiter:junit-jupiter:[5.10,6)")
    testImplementation("com.squareup.okhttp3:mockwebserver:[4.12,5)")
    testRuntimeOnly("org.junit.platform:junit-platform-launcher:[1.10,2)")
}

tasks.test {
    useJUnitPlatform()
    testLogging {
        events("failed", "skipped")
        showStandardStreams = false
    }
}

tasks.withType<JavaCompile>().configureEach {
    options.release.set(17)
    options.encoding = "UTF-8"
    options.compilerArgs.addAll(listOf("-Xlint:all", "-Werror", "-Xlint:-serial"))
}

tasks.javadoc {
    (options as StandardJavadocDocletOptions).apply {
        addStringOption("Xdoclint:none", "-quiet")
        encoding = "UTF-8"
    }
}

publishing {
    publications {
        create<MavenPublication>("maven") {
            from(components["java"])
            artifactId = "sdk"
            pom {
                name.set("Tesote Java SDK")
                description.set(project.description)
                url.set("https://www.tesote.com/docs/sdk/java")
                licenses {
                    license {
                        name.set("MIT")
                        url.set("https://opensource.org/licenses/MIT")
                    }
                }
                developers {
                    developer {
                        id.set("tesote")
                        name.set("Tesote")
                        email.set("dev.admin@tesote.com")
                    }
                }
                scm {
                    connection.set("scm:git:https://github.com/tesote/sdk.git")
                    developerConnection.set("scm:git:ssh://github.com/tesote/sdk.git")
                    url.set("https://github.com/tesote/sdk")
                }
            }
        }
    }
}

signing {
    val signingKey = providers.environmentVariable("MAVEN_GPG_KEY").orNull
    val signingPassphrase = providers.environmentVariable("MAVEN_GPG_PASSPHRASE").orNull
    if (signingKey != null && signingPassphrase != null) {
        useInMemoryPgpKeys(signingKey, signingPassphrase)
        sign(publishing.publications["maven"])
    }
}

nmcp {
    // why: nmcp 0.0.9 — `publish("<publication-name>")` block creates the task
    // `publish<PublicationName>PublicationToCentralPortal`. AUTOMATIC means
    // the Central Portal releases as soon as validation passes (no manual
    // "promote" step in the portal UI).
    publish("maven") {
        username.set(providers.environmentVariable("MAVEN_CENTRAL_USERNAME"))
        password.set(providers.environmentVariable("MAVEN_CENTRAL_PASSWORD"))
        publicationType.set("AUTOMATIC")
    }
}
