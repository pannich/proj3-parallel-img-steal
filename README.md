# Project \#3: Parallel Image Processing with Worker Stealing

### Compile
- compile
```go build -o editor editor.go```
- run bash script
```sbatch benchmark-proj3.sh```
- kill sbatch script
```
squeue --user=nichada
scancel [id]
```

Other test commands:
```go run ../editor/editor.go small
go run ../editor/editor.go small parslices 1
go run ../editor/editor.go small parslicesBSP 2
go run ../editor/editor.go small parslicesBSPOptimized 3
```

### WriteUp
- See writeup.md for report

---

## Final Project

The final project gives you the opportunity to show me what you learned
in this course and to build your own parallel system. In particular, you
should think about implementing a parallel system in the domain you are
most comfortable in (data science, machine learning, computer graphics,
etc.). The system should solve a problem that can benefit from some form
of parallelization and can be implemented in the way specified below.
I recommend reading the entire description before deciding what to implement.
If you are having trouble coming up with a problem for your system to
solve then consider the following:

-   [Embarrassingly Parallel
    Topics](https://en.wikipedia.org/wiki/Embarrassingly_parallel)
-   [Parallel
    Algorithms](https://en.wikipedia.org/wiki/Parallel_computing#Algorithmic_methods)

You are free to implement any parallel algorithm you like. However, you
are required to at least have the following features in your parallel
system:

-   An input/output component that allows the program to read in data or
    receive data in some way. The system will perform some
    computation(s) on this input and produce an output result.

-   A sequential implementation of the system. Make sure to provide a
    usage statement.

-   Basic Parallel Implementation: An implementation that uses the BSP
    pattern (using a condition variable to implement the barrier between
    supersteps), **or** a pipelining pattern (using channels) **or** a
    map-reduce implementation (again using a condition variable as barrier
    between the map and the reduce stage). Choose whichever is most suitable
    for solving the problem you have decided to tackle. The work in each
    stage or superstep should be divided among threads in a simple fashion.
    For example, if you choose an image processing problem with N images,
    then each of your T threads might be assigned to work on approximately
    N/T images. The easiest and most reasonable way to divide the work will
    depend on your problem and your chosen parallelization approach.

-   Work-stealing refinement: A work-stealing algorithm using a **dequeue**
    should be used such that the work can be split into smaller tasks, which
    are placed in a work queue such that threads will steal work from other threads
    when idle. You may either implement the dequeue as a linked-list (i.e., a chain
    of nodes similar to project \#2), or as an array as shown in class. While the
    unbounded dequeue seems more difficult to implement, the dynamic memory
    management makes it unlikely that you will suffer from the ABA problem. If you
    choose to implement the dequeue as an array, you need to ensure that a bounded
    dequeue is sufficient for your application for any valid input to your program,
    and you need to solve the ABA problem (for example using the trick of hiding a
    stamp in some bits of the integer used as array index as shown in the class
    video).

-   Provide a detailed write-up and analysis of your system. For this
    assignment, this write-up is required to have more detail to explain
    your parallel implementations since we are not giving you a problem
    to solve. See the **System Write-up** section for more details.

-   Provide all the dataset files you used in your analysis portion of
    your write up. If these files are too big then you need to provide us
    a link so we can easily download them from an external source.
    It is likely that the work-stealing refinement is only beneficial if your
    input data is structured in a certain way, e.g. if items in the input are of vastly
    different sizes, or if subtasks in your algorithm have varying or unpredictable costs.
    Make sure that this is the case for your project, so that you can showcase the pros/cons of all implementations.

-   The grade also include design points. You should think about the
    modularity of the system you are creating. Think about splitting
    your code into appropriate packages, when necessary.

-   **You must provide a script or specific commands that shows/produces
    the results of your system**. We need to be able to enter in a
    single command in the terminal window and it will run and produce
    the results of your system. Failing to provide a straight-forward
    way of executing your system that produces its result will result in
    **significant deductions** to your score. We prefer running a simple
    command line script (e.g., shell-script or python3 script). However,
    providing a few example cases of possible execution runs will be
    acceptable.

-   We should also be able to run specific versions of the system. There
    should be an option (e.g. via command line argument) to run the
    sequential version, or the various parallel versions. Please make
    sure to document this in your report or via the printing of a usage
    statement.

-   You are free to use any additional standard/third-party libraries as
    you wish. However, all the parallel work is **required** to be
    implemented by you.

-   There is a directory called `proj3` with a single `go.mod` file
    inside your repositories. Place all your work for project 3 inside
    this directory.
