module github.com/scape-labs/monorepo/shared

go 1.24

// TODO: declare the platform dependency once the platform repo's kit/ layout
// stabilises. For now we keep shared/ dependency-free so the workspace compiles
// in isolation.
//
// require github.com/scape-labs/platform v0.0.0
//
// replace github.com/scape-labs/platform => ../platform
