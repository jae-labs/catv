# Performance Improvements

This document describes the performance optimizations implemented in CATV.

## String Operations

### Issue
The `GenerateQA` function in `ollama.go` was using string concatenation with `+=` operator in a loop, which creates a new string allocation on each iteration.

### Solution
Replaced string concatenation with `bytes.Buffer.WriteString()`, which efficiently builds strings by reusing the same underlying buffer.

**Impact:** Reduces memory allocations and improves performance when processing large responses from Ollama API.

### Custom Functions Replaced
Removed custom `splitLines()` and `trimSpace()` functions in favor of standard library equivalents:
- `splitLines()` → `strings.Split()`
- `trimSpace()` → `strings.TrimSpace()` and `strings.Trim()`

**Impact:** Reduces code complexity and leverages highly optimized stdlib implementations.

## Database Operations

### Connection Pooling
Added database connection pool configuration:
- `MaxOpenConns: 25` - Maximum number of open connections
- `MaxIdleConns: 5` - Maximum number of idle connections in pool

**Impact:** Reduces overhead of creating new database connections for each operation.

### Database Indexes
Created indexes on frequently queried columns:
- `idx_flashcards_revisitin` - Index on `revisitin` column
- `idx_flashcards_file` - Index on `file` column
- `idx_flashcards_file_revisitin` - Composite index on `file` and `revisitin` columns

**Impact:** Significantly improves query performance for flashcard retrieval operations, especially when filtering by file or review status.

### Slice Pre-allocation
Pre-allocated slices with reasonable initial capacities in all query methods:
- `GetFlashcardsForReview()`: capacity 100
- `GetFlashcardsForReviewByFiles()`: capacity 50
- `GetUniqueFiles()`: capacity 20
- `GetAllFlashcards()`: capacity 100

**Impact:** Reduces memory reallocations during slice growth, improving performance when processing multiple flashcards.

## File Operations

Pre-allocated the files slice when scanning directories for markdown files with initial capacity of 10.

**Impact:** Reduces memory allocations when processing multiple markdown files.

## Benchmark Results

These optimizations provide the following benefits:
- **Reduced memory allocations**: Fewer garbage collection cycles
- **Improved query performance**: Database indexes speed up flashcard retrieval
- **Better resource utilization**: Connection pooling reduces database overhead
- **Code simplification**: Using stdlib functions reduces maintenance burden

All changes maintain backward compatibility while improving performance across the board.
