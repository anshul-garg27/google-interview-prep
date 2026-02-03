# Stage 2 - Bundle
FROM public.ecr.aws/lambda/python:3.10

RUN yum install -y gcc-c++ pkgconfig poppler-cpp-devel

COPY requirement.txt ./
COPY indic-trans-master ./
RUN pip install -r requirements.txt
RUN pip install .

RUN cp -r indictrans/models /var/lang/lib/python3.10/site-packages/indictrans/
RUN cp -r indictrans/mappings /var/lang/lib/python3.10/site-packages/indictrans/

COPY svar.csv ./
COPY vyanjan.csv ./
RUN pip install --upgrade pip && \
    pip install -r requirement.txt &&\
    pip install --upgrade numpy

COPY fake.py ./
COPY baby_names_.csv ./baby_names.csv
RUN rm -r indictrans
CMD [ "fake.handler" ]